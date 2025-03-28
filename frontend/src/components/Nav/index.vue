<template>
    <div class="Nav">
        <PayloadDialog />
        <el-segmented v-model="navStore.nav" :options="options" size="large" @change="onNavChange" />
        <Card class=card ref="card" />
    </div>
</template>

<script setup lang='ts'>
import Card from './Card.vue';
import { ref, watch } from 'vue'
import { useNavStore } from '../../store/nav'
import PayloadDialog from './PayloadDialog.vue';

const navStore = useNavStore()

const card = ref(null)
const options = [
    'Reverse',
    'Bind',
    'MSFVenom',
    'HoaxShell',
]

function onNavChange() {
    card.value.updateList()
}

</script>

<style lang="scss" scoped>
.Nav {
    // width: 100%;
    padding: 0 52px;
}
.card{
    // max-height: 480px;
    // overflow-y: auto;
}
</style>